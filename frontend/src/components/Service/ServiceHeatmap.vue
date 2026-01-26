<template>
    <apexchart v-if="ready" width="100%" height="180" type="heatmap" :options="plotOptions" :series="series" @dataPointSelection="onDataPointClick"></apexchart>
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
          await this.chartHeatmap()
      },
      data() {
          return {
              ready: false,
              data: [],
              dateMap: {}, // Map to store date info: key = "seriesIndex-dataPointIndex", value = date string
              plotOptions: {
                  chart: {
                      selection: {
                          enabled: false
                      },
                      zoom: {
                          enabled: false
                      },
                      toolbar: {
                          show: false
                      },
                      events: {
                          dataPointSelection: (event, chartContext, config) => {
                              // This will be handled by the @dataPointSelection event
                          }
                      }
                    },
                  theme: {
                    mode: "dark",
                    },
                      colors: [ "#cb3d36" ],
                      enableShades: true,
                      shadeIntensity: 0.5,
                      colorScale: {
                          ranges: [ {
                              from: -1,
                              to: -1,
                              color: '#666666',
                              name: 'future',
                          },
                          {
                              from: 0,
                              to: 0,
                              color: '#28a745',
                              name: 'none',
                          },
                              {
                                  from: 1,
                                  to: 5,
                                  color: '#ffc107',
                                  name: 'low',
                              },
                              {
                                  from: 6,
                                  to: 10,
                                  color: '#fd7e14',
                                  name: 'medium',
                              },
                              {
                                  from: 11,
                                  to: 999,
                                  color: '#dc3545',
                                  name: 'high',
                              }
                          ]
                      },
                      dataLabels: {
                          enabled: true,
                          formatter: function(val) {
                              if (val === -1 || val < 0) {
                                  return 'F';
                              }
                              return val;
                          },
                          style: {
                              colors: ['#ffffff']
                          }
                      },
                  xaxis: {
                      tickAmount: '1',
                      tickPlacement: 'between',
                      min: 1,
                      max: 31,
                      type: "numeric",
                      labels: {
                          show: true
                      },
                      tooltip: {
                          enabled: false
                      }
                  },
                  yaxis: {
                      labels: {
                          show: true
                      },
                  },
                  tooltip: {
                      theme: false,
                      enabled: true,
                      custom: ({series, seriesIndex, dataPointIndex, w}) => {
                          const dateKey = `${seriesIndex}-${dataPointIndex}`;
                          const dateStr = this.dateMap[dateKey];
                          const dataPoint = w.globals.initialSeries[seriesIndex].data[dataPointIndex];
                          const failures = dataPoint ? (dataPoint.y || 0) : 0;
                          const isFuture = dataPoint && dataPoint.y === -1;
                          
                          if (dateStr) {
                              const date = this.parseISO(dateStr);
                              const formattedDate = date.toLocaleDateString('en-us', { 
                                  weekday: 'long', 
                                  year: 'numeric', 
                                  month: 'long', 
                                  day: 'numeric' 
                              });
                              if (isFuture) {
                                  return `<div style="padding: 12px 16px 12px 16px; background-color: rgba(0, 0, 0, 0.9); border-radius: 4px; color: #ffffff;">Future Date<br>${formattedDate}</div>`;
                              }
                              return `<div style="padding: 12px 16px 12px 16px; background-color: rgba(0, 0, 0, 0.9); border-radius: 4px; color: #ffffff;">${failures} ${failures === 1 ? 'Failure' : 'Failures'}<br>${formattedDate}</div>`;
                          }
                          
                          // Fallback
                          return `<div style="padding: 12px 16px 12px 16px; background-color: rgba(0, 0, 0, 0.9); border-radius: 4px; color: #ffffff;">${failures} ${failures === 1 ? 'Failure' : 'Failures'}</div>`;
                      },
                      fixed: {
                          enabled: true,
                          position: 'topLeft',
                          offsetX: 0,
                          offsetY: 0,
                      }
                  }
                  },
              series: [ {
                  data: []
              } ]
          }
      },
      methods: {
          async chartHeatmap() {
            const monthData = []
            let start = this.firstDayOfMonth(this.now())
            this.dateMap = {} // Reset date map

            for (let i=0; i<3; i++) {
                const monthInfo = await this.heatmapData(this.addMonths(start, -i), this.lastDayOfMonth(this.addMonths(start, -i)), i)
                monthData.push(monthInfo)
            }

            this.series = monthData
            this.ready = true
          },
          async heatmapData(start, end, seriesIndex) {
              const data = await Api.service_failures_data(this.service.id, this.toUnix(start), this.toUnix(end), "24h", true)
              let dataArr = []
              const now = this.now();
              const today = this.beginningOf('day', now);
              
              // Create a map of existing data by day of month
              const dataMap = {};
              if (data) {
                  data.forEach((d) => {
                      const date = this.parseISO(d.timeframe);
                      const dayOfMonth = date.getDate();
                      dataMap[dayOfMonth] = {
                          amount: d.amount,
                          date: d.timeframe
                      };
                  });
              }
              
              // Fill in all days of the month (1-31)
              const daysInMonth = this.lastDayOfMonth(start).getDate();
              let dataIndex = 0;
              for (let day = 1; day <= daysInMonth; day++) {
                  const testDate = new Date(start.getFullYear(), start.getMonth(), day);
                  const dateStartOfDay = this.beginningOf('day', testDate);
                  const isFuture = dateStartOfDay > today;
                  
                  if (dataMap[day]) {
                      // We have data for this day
                      const dateKey = `${seriesIndex}-${dataIndex}`;
                      this.dateMap[dateKey] = dataMap[day].date;
                      dataArr.push({
                          x: day,
                          y: isFuture ? -1 : dataMap[day].amount,
                          isFuture: isFuture
                      });
                      dataIndex++;
                  } else {
                      // No data for this day - fill it in
                      const dateKey = `${seriesIndex}-${dataIndex}`;
                      const dateStr = testDate.toISOString();
                      this.dateMap[dateKey] = dateStr;
                      dataArr.push({
                          x: day,
                          y: isFuture ? -1 : 0,
                          isFuture: isFuture
                      });
                      dataIndex++;
                  }
              }

              return {
                  name: start.toLocaleString('en-us', { month: 'long'}), 
                  data: dataArr,
                  monthStart: start
              }
          },
          onDataPointClick(event, chartContext, config) {
              if (config.dataPointIndex !== undefined && config.seriesIndex !== undefined) {
                  const dateKey = `${config.seriesIndex}-${config.dataPointIndex}`;
                  const dateStr = this.dateMap[dateKey];
                  
                  if (dateStr) {
                      const date = this.parseISO(dateStr);
                      const dateStartOfDay = this.beginningOf('day', date);
                      const today = this.beginningOf('day', this.now());
                      
                      // Don't allow clicking on future days
                      if (dateStartOfDay > today) {
                          return;
                      }
                      
                      const series = this.series[config.seriesIndex];
                      this.$emit('day-selected', {
                          date: date,
                          month: series.name
                      });
                  }
              }
          }
      }
  }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
