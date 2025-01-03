Feature: display the current offline status

  Background:
    Given a Git repo with origin

  Scenario: default value
    When I run "git-town offline"
    Then Git Town prints:
      """
      no
      """

  Scenario Outline: configured in local Git metadata
    Given global Git setting "git-town.offline" is "<VALUE>"
    When I run "git-town offline"
    Then Git Town prints:
      """
      <OUTPUT>
      """
    Examples:
      | VALUE | OUTPUT |
      | yes   | yes    |
      | on    | yes    |
      | true  | yes    |
      | 1     | yes    |
      | t     | yes    |
      | no    | no     |
      | off   | no     |
      | false | no     |
      | f     | no     |
      | 0     | no     |

  Scenario: invalid value in Git metadata
    Given global Git setting "git-town.offline" is "zonk"
    When I run "git-town offline"
    Then Git Town prints the error:
      """
      invalid value for git-town.offline: "zonk". Please provide either "yes" or "no"
      """
