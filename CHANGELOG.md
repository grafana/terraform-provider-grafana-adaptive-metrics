## Unreleased

## v0.4.0

- [BREAKING] Default omitted auto-apply configuration to disabled. Previously an omitted `auto_apply` block preserved the server-side value; it is now planned as `enabled = false`, so set it explicitly on any segment or recommendations config where auto-apply should stay on.
- [CHANGE] Remove the `grafana-adaptive-metrics_policy` resource and the `policy_id` field from the segment resource. The underlying feature was never made available, so no existing configuration should reference them.
- [FEATURE] Add auto-apply gate policy support to segment and recommendations config resources.

## v0.3.6

- [BUGFIX] Stop overriding a segment's `policy_id` with the default policy when the field is removed from the configuration.

## v0.3.5

- [CHANGE] Drop the private preview notice from the `auto_apply` attribute of the segment and recommendations config resources.
- [BUGFIX] Mark the segment `auto_apply` attribute as computed so a server-side value no longer causes a perpetual diff.

## v0.3.4

- [ENHANCEMENT] Document how the `ruleset` and `rule` resources conflict when both manage the same metric.

## v0.3.3

- [FEATURE] Add policy resource.
- [ENHANCEMENT] Add an ability to assign policy to a segment using `policy_id` field in `segment` resource.

## v0.3.2

- [FEATURE] Add `auto_apply` attribute to the recommendations config resource (private preview).

## v0.3.1

- [FEATURE] Add `auto_apply` attribute to the segment resource to configure server-side auto-apply (private preview).
- [BUGFIX] Fix the permissions of the release workflow.

## v0.3.0

- [FEATURE] Add segment resource
- [FEATURE] Add ruleset resource
- [ENHANCEMENT] Add `segment` selector to recommendations datasource
- [ENHANCEMENT] Add `segment` field to exemption resource
- [ENHANCEMENT] Make user-agent header configurable

## v0.2.0

- [ENHANCEMENT] Add `reason`, `disable_recommendations`  field to exemptions resource

## v0.1.2

- [ENHANCEMENT] Add `auto_import` field to rules resource to toggle create vs. create/update behavior

## v0.1.1

- [ENHANCEMENT] Minor doc updates

## v0.1.0

Inital release.

- [FEATURE] Add provider resource
- [FEATURE] Add rule resource
- [FEATURE] Add exemptions resource
- [FEATURE] Add recommendations datasource
