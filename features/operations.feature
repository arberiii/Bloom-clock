# file: $GOPATH/bloom-clock/features/operations.feature
Feature: Operations
  In order to prove the happened before operation
  As a User
  I need to send requests to nodes


Scenario: false error
    Given two timestamps
    When requesting the happened-before operation
    And both timestamps are comparable
    And the first provided element is smaller than the second provided element
    Then return the inferred false positive rate associated and a false error

Scenario: true error
    Given two timestamps
    When requesting the happened-before operation
    And both timestamps are comparable
    And the first provided element is larger than the second provided element
    Then return zero and a true error

Scenario: incomparable timestamps
    Given two timestamps
    When requesting the happened-before operation
    And the timestamps are not comparable
    Then return zero and a true error

Scenario: false error
    Given two timestamps
    When requesting the happened-after operation
    And both timestamps are comparable
#    And the first provided element is larger than the second provided element
    Then return the inferred false positive rate and a false error

Scenario: happened after error
    Given two timestamps
    When requesting the happened-after operation
    And both timestamps are comparable
    And the first provided element is smaller than the second element
    Then return one and a true error

Scenario: incomparable
    Given two timestamps
    When requesting the happened-after operation
    And the timestamps are not comparable
    Then return zero and a true error

Scenario: comparable
    Given two timestamps
    When requesting the compare operation
    And the two timestamps are comparable
    Then return true

Scenario: timestamps are not comparable
    Given two timestamps
    And the two timestamps are not comparable
    When requesting the compare operation
    Then return false

Scenario: isLarger
    Given two timestamps
    And the first timestamp is larger than the second provided one
    When requesting the isLarger operation
    And the two timestamps are comparable
    Then return that is larger

Scenario: isLarger false
    Given two timestamps
    And the first timestamp is smaller than the second provided one
    When requesting the isLarger operation
    And the two timestamps are comparable
    Then return that is not larger

Scenario: overlap
    Given two timestamps
    And the first timestamp is larger than the second provided one
    When requesting the overlaps operation
    Then return that is comparable and larger

Scenario: overlap but false
    Given two timestamps
    And the first timestamp is smaller than the second provided one
    When requesting the overlaps operation
    Then return that is not comparable or larger




