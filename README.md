# Uli

An easy-to-use command line tool for fetching real-time sports information 

## Motivation

I recently switch to DuckDuckGo for my internet searches. One thing I miss about good is it's ability to show me today's sports games. Uli solves this problem for me by putting that information into my terminal.

# Usage

    $ uli nhl
    +---------+---------------------+----------------------+----+-------+
    |  TIME   |        HOME         |         AWAY         |    | SCORE |
    +---------+---------------------+----------------------+----+-------+
    | 3:00 pm | St. Louis Blues     | Minnesota Wild       | 🏁 | 2 - 3 |
    | 5:00 pm | Florida Panthers    | Ottawa Senators      | 🏁 | 5 - 1 |
    | 5:00 pm | Washington Capitals | Arizona Coyotes      | 🏁 | 1 - 4 |
    | 7:00 pm | Winnipeg Jets       | New Jersey Devils    | 🔴 | 3 - 2 |
    | 7:00 pm | Boston Bruins       | Vegas Golden Knights | 🔴 | 3 - 0 |
    | 9:00 pm | San Jose Sharks     | Calgary Flames       |    |       |
    | 9:30 pm | Edmonton Oilers     | Colorado Avalanche   |    |       |
    +---------+---------------------+----------------------+----+-------+

    $ uli nhl nyr
    +-------------------+----+-----------------------+--------------------+-------+
    |       TIME        |    |         HOME          |        AWAY        | SCORE |
    +-------------------+----+-----------------------+--------------------+-------+
    | 11/04 7:00 pm     | 🏁 | New York Rangers      | Buffalo Sabres     | 3 - 1 |
    | 11/06 7:00 pm     | 🏁 | New York Rangers      | Montréal Canadiens | 5 - 3 |
    | 11/09 7:30 pm     | 🏁 | Detroit Red Wings     | New York Rangers   | 3 - 2 |
    | Yesterday 7:00 pm | 🏁 | Columbus Blue Jackets | New York Rangers   | 4 - 5 |
    | Tomorrow 7:00 pm  |    | New York Rangers      | Vancouver Canucks  |       |
    | 11/15 7:00 pm     |    | New York Islanders    | New York Rangers   |       |
    | 11/17 7:00 pm     |    | New York Rangers      | Florida Panthers   |       |
    +-------------------+----+-----------------------+--------------------+-------+

You can also alias Uli commands to make them even shorter by adding lines like this to your bashrc or zshrc files:

    alias nhl='uli nhl'
    alias nyr='uli nhl nyr'