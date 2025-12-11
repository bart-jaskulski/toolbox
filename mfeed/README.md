# mfeed – Meta Feed Generator

The Meta Feed Generator is a specialized tool designed to streamline the consumption of content from multiple RSS feeds. By aggregating feeds, ranking their contents based on personalized criteria, and generating a curated Atom feed, this tool ensures that users have access to the best articles tailored to their interests without the need to manually sift through numerous entries. This solution is particularly beneficial for avid readers and researchers who rely on staying updated with high-quality content from various sources.

## Features

- **Feed Aggregation**: Combines multiple RSS feeds into a single source of content, making it easier to manage and consume information.
- **Content Ranking**: Articles are dynamically ranked based on user-defined preferences, such as topic relevance, article depth, and quality of insights. This ensures that the feed prioritizes content that aligns with the user’s interests.
- **Atom Feed Generation**: Outputs a refined list of top articles in an Atom feed format, which can be easily integrated with feed readers or websites.
- **Personalization**: Users can customize ranking criteria and the AI model used for analyzing articles, allowing for a highly personalized content curation process.
- **OpenAI API Compatible LLM Integration**: Utilizes OpenAI’s API to intelligently assess and score articles, ensuring that the feed includes only the most engaging and relevant content.

## Installation

Install the CLI:
```bash
go install github.com/bart-jaskulski/toolbox/mfeed@latest
```

Configure your RSS sources by listing feed URLs in the `feeds` file in the project root (or wherever you run `mfeed` from).

## Usage

Execute the built application to process the feeds. The output will be an Atom feed consisting of the top-ranked articles, printed to the standard output. This output can be redirected to a file or directly used by feed readers to display the curated content list.

## Configuration

Adjustments to AI prompts and model choices can be made in the `./prompts/` directory and within the `./main.go` file for re-ranking based on new preferences. These modifications require recompilation to take effect.

## Contributing

Contributors are encouraged to submit pull requests or open issues to propose features or report bugs. For significant changes, please discuss them in an issue first.

## License

This project is open-sourced under the MIT License. See the `LICENSE` file for more details.
