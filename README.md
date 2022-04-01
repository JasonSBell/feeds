# Feeds

This repository contains the code for Allokate's Feeds service. This service contains a series of scripts executed as cron jobs that scrape data from various sources (RSS feeds, Nasdaq's API, etc). It also contains the API for serving out some of this feed driven data. It also provides endpoints for fetching recent activity for a basket of stocks.

# Table of Contents

- [Feeds](#feeds)
- [Table of Contents](#table-of-contents)
- [Endpoints](#endpoints)

# Endpoints

- GET /api/ping
- GET /api/earnings
- GET /api/dividends
- GET /api/news
- GET /api/tickers/:ticker
- GET /api/feeds/activity
