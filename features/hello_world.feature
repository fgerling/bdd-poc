Feature: Hello World

  Scenario: Say Hello
    Given I have "echo" in PATH
    When I run "echo Hello World"
    Then the output contains "Hello World"
