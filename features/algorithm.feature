# file: $GOPATH/bloom-clock/features/algorithm.feature
Feature: Algorithm
 In order to compare timestamps between nodes
 As a User
 I need to send requests to nodes

Scenario: latest timestamp
    Given network topology
    When the node receives a comparable timestamp that happens-before the node’s timestamp
    Then assign as the latest timestamp the node’s original timestamp

Scenario: max timestamp
    Given network topology
    When the node receives a comparable timestamp that happens-after A’s timestamp
    Then assign as the latest timestamp the node’s receiving timestamp

Scenario: conflict merge
    Given network topology
    When the node receives an incomparable timestamp,
    Then assign as the latest timestamp a bloom filter that contains the max element of both bloom filters

Scenario: broadcast element
    Given network topology
    When the node has an internal event,
    Then hash the event / element
    And increment the internal bloom filter
    And send it to all the nodes of the network

Scenario: run a simple case of algorithm
    Given network topology
    When the node receives a bloom filter
    Then update the node’s latest timestamp


