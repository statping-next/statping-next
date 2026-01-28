<template>
    <div>
        <div v-if="!ready" class="text-center py-4">
            <font-awesome-icon icon="circle-notch" spin/>
            <span class="ml-2">Loading...</span>
        </div>

        <div v-if="ready" style="margin-bottom: 1rem;">
            <table style="width: 100%; border-collapse: collapse; font-size: 0.65rem; table-layout: fixed;">
                <thead>
                    <tr>
                        <th style="background-color: rgba(0, 0, 0, 0.3); font-weight: bold; border: 1px solid rgba(255, 255, 255, 0.1); padding: 0.25rem; text-align: center; width: 80px;"></th>
                        <th v-for="day in 31" :key="day" style="background-color: rgba(0, 0, 0, 0.2); font-weight: bold; font-size: 0.6rem; border: 1px solid rgba(255, 255, 255, 0.1); padding: 0.1rem; text-align: center;">{{ day }}</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(monthData, monthIndex) in monthsData" :key="monthIndex">
                        <td style="background-color: rgba(0, 0, 0, 0.3); font-weight: bold; text-align: left; padding-left: 0.5rem; border: 1px solid rgba(255, 255, 255, 0.1); padding: 0.25rem; font-size: 0.65rem;">{{ monthData.monthName }}</td>
                        <td 
                            v-for="(dayData, dayIndex) in monthData.days" 
                            :key="dayIndex"
                            :style="{
                                backgroundColor: dayData.color,
                                color: '#ffffff',
                                fontWeight: '500',
                                cursor: (dayData.isFuture || !canViewFailures) ? 'not-allowed' : 'pointer',
                                opacity: dayData.isFuture ? 0.6 : 1,
                                border: '1px solid rgba(255, 255, 255, 0.1)',
                                padding: '0.1rem',
                                textAlign: 'center',
                                fontSize: '0.65rem'
                            }"
                            @click="!dayData.isFuture && canViewFailures && onDayClick(dayData)"
                            :title="dayData.isFuture ? dayData.tooltip : (canViewFailures ? dayData.tooltip : 'Login required to view failure details')"
                        >
                            {{ dayData.displayValue }}
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Failure list popout -->
        <div v-if="selectedDayFailures.length > 0 && canViewFailures" class="mt-4">
            <div class="card">
                <div class="card-header d-flex justify-content-between align-items-center">
                    <h5 class="mb-0">Failures on {{ selectedDayDate }}</h5>
                    <button class="btn btn-sm btn-secondary" @click="closeFailureList">×</button>
                </div>
                <div class="card-body">
                    <div v-if="selectedDayFailures.length === 0" class="text-muted">
                        No failures found for this day.
                    </div>
                    <table v-else class="table">
                        <thead>
                            <tr>
                                <th scope="col">#</th>
                                <th scope="col">Issue</th>
                                <th scope="col">Status Code</th>
                                <th scope="col">Ping</th>
                                <th scope="col">Created</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="(failure, index) in selectedDayFailures" :key="index">
                                <th class="font-1" scope="row">{{failure.id}}</th>
                                <td class="font-1">{{failure.issue}}</td>
                                <td class="font-1">{{failure.error_code}}</td>
                                <td class="font-1">{{humanTime(failure.ping || failure.ping_time)}}</td>
                                <td class="font-1">{{ago(failure.created_at)}}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
  import Api from "../../API"

  export default {
      name: 'ServiceHeatmap',
      props: {
          service: {
              type: Object,
              required: true
          }
      },
      async created() {
          await this.loadData()
      },
      data() {
          return {
              ready: false,
              monthsData: [],
              selectedDayFailures: [],
              selectedDayDate: ''
          }
      },
      computed: {
          canViewFailures() {
              // Check if user is logged in
              const isLoggedIn = this.$store.getters.loggedIn || this.$store.getters.user || this.$store.getters.admin
              // Check if setting allows unauthenticated users to view failures
              const core = this.$store.getters.core
              const allowUnauthenticated = core && core.show_failures_to_unauthenticated
              return isLoggedIn || allowUnauthenticated
          }
      },
      methods: {
          async loadData() {
              const monthsData = []
              let start = this.firstDayOfMonth(this.now())
              const now = this.now()
              const today = this.beginningOf('day', now)

              for (let i = 0; i < 3; i++) {
                  const monthStart = this.addMonths(start, -i)
                  const monthEnd = this.lastDayOfMonth(monthStart)
                  const monthInfo = await this.loadMonthData(monthStart, monthEnd, today)
                  monthsData.push(monthInfo)
              }

              // Reverse so oldest month is first
              this.monthsData = monthsData.reverse()
              this.ready = true
          },
          async loadMonthData(monthStart, monthEnd, today) {
              // Fetch individual failures instead of pre-aggregated data
              // This ensures consistency between display and click handler
              const failures = await Api.service_failures(
                  this.service.id,
                  this.toUnix(monthStart),
                  this.toUnix(monthEnd),
                  999999,
                  0
              )

              // Group failures by day of month in local timezone
              const dataMap = {}
              if (failures && failures.length > 0) {
                  failures.forEach((failure) => {
                      const date = this.parseISO(failure.created_at)
                      const dayOfMonth = date.getDate()

                      if (!dataMap[dayOfMonth]) {
                          dataMap[dayOfMonth] = {
                              amount: 0,
                              date: date.toISOString(),
                              dateObj: date,
                              failures: []
                          }
                      }
                      dataMap[dayOfMonth].amount++
                      dataMap[dayOfMonth].failures.push(failure)
                  })
              }

              // Build days array for this month
              const daysInMonth = this.lastDayOfMonth(monthStart).getDate()
              const days = []

              for (let day = 1; day <= daysInMonth; day++) {
                  // Create date in local timezone for display
                  const testDate = new Date(monthStart.getFullYear(), monthStart.getMonth(), day)
                  const dateStartOfDay = this.beginningOf('day', testDate)
                  const isFuture = dateStartOfDay > today

                  let value, dateStr, dateObj

                  if (dataMap[day]) {
                      value = isFuture ? -1 : dataMap[day].amount
                      dateStr = dataMap[day].date
                      dateObj = dataMap[day].dateObj
                  } else {
                      value = isFuture ? -1 : 0
                      dateStr = testDate.toISOString()
                      dateObj = testDate
                  }

                  const color = this.getColorForValue(value)
                  const displayValue = value === -1 ? 'F' : value

                  const formattedDate = dateObj.toLocaleDateString('en-us', {
                      weekday: 'long',
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric'
                  })

                  const tooltip = isFuture 
                      ? `Future Date: ${formattedDate}`
                      : `${value} ${value === 1 ? 'Failure' : 'Failures'}: ${formattedDate}`

                  days.push({
                      day: day,
                      value: value,
                      color: color,
                      displayValue: displayValue,
                      isFuture: isFuture,
                      dateStr: dateStr,
                      dateObj: dateObj,
                      tooltip: tooltip
                  })
              }

              // Fill remaining days up to 31 with empty cells
              for (let day = daysInMonth + 1; day <= 31; day++) {
                  days.push({
                      day: day,
                      value: null,
                      color: 'transparent',
                      displayValue: '',
                      isFuture: false,
                      dateStr: null,
                      dateObj: null,
                      tooltip: ''
                  })
              }

              return {
                  monthName: monthStart.toLocaleString('en-us', { month: 'long' }),
                  monthStart: monthStart,
                  days: days
              }
          },
          getColorForValue(value) {
              if (value === null || value === undefined) {
                  return 'transparent'
              }

              const numValue = Number(value)
              if (isNaN(numValue)) {
                  return '#28a745' // default green
              }

              if (numValue === -1) {
                  return '#666666' // grey for future
              } else if (numValue === 0) {
                  return '#28a745' // green for no failures
              } else if (numValue >= 1 && numValue <= 5) {
                  return '#ffc107' // yellow for low
              } else if (numValue >= 6 && numValue <= 10) {
                  return '#fd7e14' // orange for medium
              } else if (numValue >= 11) {
                  return '#dc3545' // red for high
              }

              return '#28a745' // default green
          },
          async onDayClick(dayData) {
              // Check if user can view failures
              if (!this.canViewFailures) {
                  return
              }
              if (!dayData.dateStr || dayData.isFuture) {
                  return
              }

              try {
                  // Get the start and end of the selected day in local timezone
                  const dayStart = this.beginningOf('day', dayData.dateObj)
                  const dayEnd = this.endOf('day', dayData.dateObj)

                  // Fetch failures for this specific day
                  const failures = await Api.service_failures(
                      this.service.id,
                      this.toUnix(dayStart),
                      this.toUnix(dayEnd)
                  )

                  this.selectedDayFailures = failures || []
                  this.selectedDayDate = dayData.dateObj.toLocaleDateString('en-us', {
                      weekday: 'long',
                      year: 'numeric',
                      month: 'long',
                      day: 'numeric'
                  })
              } catch (error) {
                  console.error('Error loading failures:', error)
                  this.selectedDayFailures = []
              }
          },
          closeFailureList() {
              this.selectedDayFailures = []
              this.selectedDayDate = ''
          }
      }
  }
</script>

<style scoped>
</style>
