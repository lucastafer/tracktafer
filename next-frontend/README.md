# TrackTafer - Vehicle tracking system & freight calculator - FrontEnd

## Description

Simple tracking system's front-end, built with the latest version of Next.js, using SSR and integrated with Google Maps and our API.

## Requirements

It's mandatory integrate with Google Maps API, for that, you need to generate your Google Cloud API Key. Follow these steps:

1. Access [Google Cloud Platform](https://cloud.google.com/).
2. Create new project.
3. Activate the following API's: **Places**, **Directions** and **Maps Javascript (from Google Maps)**.
4. Save the generated API key (this will be your NEXT_PUBLIC_GOOGLE_MAPS_API_KEY env at your .env file).

**Attention: Remember not to leave your Google Maps API Key public (or deactivate it), as this could lead to unwanted charges.**

## Run application

To run only the front-end without using Docker:

Generate your `.env` with:

```
cp .env.example .env
```

Install dependencies:

```bash
npm install
```

Run server:

```bash
npm run dev
```
