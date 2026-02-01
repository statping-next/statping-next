// Default branding
// Empty 1x1 transparent favicon (no flash); used in HTML before app loads
export const EMPTY_FAVICON = '/img/empty-favicon.gif'
// Project default favicon when user has not set a custom one
export const DEFAULT_FAVICON = '/img/default-favicon.ico'
// Default logo when admin has not set a custom one (never shown if user set a logo)
export const DEFAULT_LOGO = '/img/default-logo.svg'

/** Returns the logo URL if set (non-empty), otherwise null so template can use DEFAULT_LOGO. */
export function logoUrlIfSet (value) {
  if (value == null) return null
  const s = String(value).trim()
  return s.length > 0 ? s : null
}
