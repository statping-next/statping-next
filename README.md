# Statping NEXT - Status Page & Monitoring Server

![Logo Image](/frontend/public/img/default-logo.svg)

Statping NEXT is a drop-in replacement for [statping-ng](https://github.com/statping-ng/statping-ng). It was conceived after the author desired several changes to the project, that were likely not ever to be implemented, both due to their opinionated nature and due to a seeming stall in that project's development.

Statping NEXT features several significant changes from the parent project, including:

* A more compact, stylish user-focused main page that clearly shows the current service status, check/down times, announcements, and incidents, without additional information like graphs.
* Light and dark theming (persistent in browser storage, default configurable by the admin).
* Automatic page refresh at 15 and 60 second intervals (default configurable by the admin), in addition to the traditional one-and-done load.
* Improved navigation buttons.
* Customized favicon and logo image support, with both light and dark images supported, and a nice following navbar with your logo always visible during page scrolling (SVG highly recommended).
* Significantly improved main-page loading time over `statping-ng` due to eliminating most of the heavy DB queries that were required for the original home page to show its multiple graphs.
* An improved per-service heatmap with per-day failure reporting (private by default, optionally public).
* Significantly improved and streamlined incident support, including global incidents (in addition to per-service incidents), archiving of incidents, and auto-archiving on resolution.
* Rearrangement and reorganization of the admin dashboard, specifically in the top links - important day-to-day features are listed first, and incidents and announcements are consolidated in the "Incidents" page. The default dashboard page has also been renamed/moved to "Status".

Statping NEXT is currently providing monitoring for [the Jellyfin Project](https://status.jellyfin.org) and the author's personal systems.

Note: Your Statping NEXT service should not be running on the same instance you're trying to monitor. If your server crashes your Status Page should still remaining online to notify your users of downtime.

## Still (Mostly) Compatible with `statping-ng`

Most of the [documentation for statping-ng](https://github.com/statping-ng/statping-ng/wiki) still applies fully to Statping NEXT. We've strived not to change too much core functionality or break compatibility in either the backend or the API, though there are changes to support the new features mentioned above. The build process also still works the exact same way, though at this time we are not providing binaries. Use `make binary` to build a Linux binary, or reference the linked Wiki for instructions on how to build for other platforms.

## Screenshots

### Dark Theme

![Dark Theme - Jellyfin](/images/dark-theme-jellyfin.png)

### Light Theme

![Light Theme - Jellyfin](/images/light-theme-jellyfin.png)

# Contributing

Statping NEXT accepts Push Requests to the `master` branch!

This is an **LLM-friendly** project. Much of the development beyond [statping-ng](https://github.com/statping-ng/statping-ng) was completed with the assistance of LLMs. If this makes you uncomfortable, you may use the parent project instead. We still expect any PRs to **respect formatting and code cleanliness standards**.

Versioning will begin with `v1.0` once the initial development of new features has concluded, and from there will increment in a major-minor format; major releases **may** introduce API changes, while minor releases will only be API-compatible bugfixes.
