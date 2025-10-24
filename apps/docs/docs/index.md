---
sidebar_position: 1
slug: /
---

# Formbricks Store

The foundational data service for the Formbricks experience management platform.

## What is Store?

Store is an open-source, self-hostable microservice for collecting and serving experience data – survey responses, NPS scores, product reviews, support feedback, and more. Built in Go for performance and simplicity.

**Part of the Formbricks Ecosystem**: Store serves as the data foundation for survey delivery, analytics, AI agents, and custom workflows.

## Why Store?

- 🚀 **High Performance**: Go-powered microservice handles high write volume
- 🤖 **AI-Powered**: Automatic sentiment analysis, topic extraction, and semantic search
- 📊 **Analytics-Ready**: Optimized for Apache Superset, Power BI, Tableau, Looker, Snowflake
- 🔐 **Self-Hostable**: Run in a single Docker container with PostgreSQL
- 🛠️ **Developer-Friendly**: OpenAPI spec, webhooks, clean REST API

## Use Cases

### Centralize Feedback Data

Collect experience data from multiple sources into one unified data store:
- Survey responses from your app or website
- App Store and Google Play reviews
- Trustpilot and other review platforms
- Support ticket feedback (Intercom, Zendesk)
- Social media sentiment

### Power Analytics Dashboards

Store's schema is optimized for direct SQL queries and BI tool integration:
- Connect Apache Superset for real-time dashboards
- Build Power BI or Tableau reports without complex ETL
- Export to Snowflake or Redshift for data warehousing
- Perform time-series analysis on feedback trends

### AI-Powered Insights

Store automatically enriches text feedback with actionable insights:
- **Sentiment analysis**: Positive, negative, or neutral classification
- **Emotion detection**: Joy, anger, sadness, and more
- **Topic extraction**: Automatically identify themes in feedback
- **Semantic search**: Find similar responses using natural language queries

[Learn more about AI enrichment →](./core-concepts/ai-enrichment)

### Build Custom Workflows

Use webhooks to trigger actions based on enriched feedback:
- Send Slack notifications for low NPS scores or negative sentiment
- Route negative reviews to your support team automatically
- Update CRM records with AI-extracted insights
- Build custom dashboards with semantic search capabilities

### Future: Connector Ecosystem

An open connector ecosystem is planned to simplify data integration. [Learn more →](./core-concepts/connectors)

## Key Features

### AI Enrichment
Automatic sentiment, emotion, and topic analysis powered by OpenAI for text feedback.

### Semantic Search
Find similar feedback using natural language queries with pgvector embeddings.

### Real-Time Webhooks
Trigger workflows instantly when new feedback arrives or insights are generated.

### SQL-Friendly Schema
Query directly with your favorite BI tool – no complex transformations needed.

### Open Source
Apache 2.0 licensed. Self-host, modify, and own your data completely.

## Quick Links

<div className="row">
  <div className="col col--6">
    <div className="card margin-bottom--lg">
      <div className="card__header">
        <h3>🚀 Get Started</h3>
      </div>
      <div className="card__body">
        <p>Set up Store locally in 5 minutes</p>
      </div>
      <div className="card__footer">
        <a href="./quickstart" className="button button--primary button--block">
          Quick Start Guide
        </a>
      </div>
    </div>
  </div>
  <div className="col col--6">
    <div className="card margin-bottom--lg">
      <div className="card__header">
        <h3>📚 Learn the Basics</h3>
      </div>
      <div className="card__body">
        <p>Understand Store's data model and concepts</p>
      </div>
      <div className="card__footer">
        <a href="./core-concepts/data-model" className="button button--secondary button--block">
          Core Concepts
        </a>
      </div>
    </div>
  </div>
</div>

<div className="row">
  <div className="col col--6">
    <div className="card margin-bottom--lg">
      <div className="card__header">
        <h3>🔧 API Reference</h3>
      </div>
      <div className="card__body">
        <p>Interactive API documentation and testing</p>
      </div>
      <div className="card__footer">
        <a href="./api-reference" className="button button--secondary button--block">
          Explore API
        </a>
      </div>
    </div>
  </div>
  <div className="col col--6">
    <div className="card margin-bottom--lg">
      <div className="card__header">
        <h3>⚙️ Configuration</h3>
      </div>
      <div className="card__body">
        <p>Environment variables and settings reference</p>
      </div>
      <div className="card__footer">
        <a href="./reference/environment-variables" className="button button--secondary button--block">
          Configuration Guide
        </a>
      </div>
    </div>
  </div>
</div>

## Community & Support

- **GitHub**: [formbricks/store](https://github.com/formbricks/store)
- **Discussions**: [Ask questions and share ideas](https://github.com/formbricks/store/discussions)
- **Issues**: [Report bugs or request features](https://github.com/formbricks/store/issues)
- **Documentation**: You're reading it! 📖
