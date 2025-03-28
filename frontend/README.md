This is a [Next.js](https://nextjs.org) project using the Pages Router architecture.

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `pages/index.tsx`. The page auto-updates as you edit the file.

## Project Structure

```
frontend/
  ├── pages/             # Pages Router directory
  │   ├── _app.tsx      # Root component
  │   ├── _document.tsx # Custom document
  │   └── index.tsx     # Main page
  ├── styles/           # Global styles
  ├── components/       # Reusable components
  ├── hooks/           # Custom React hooks
  ├── utils/           # Utility functions
  ├── types/           # TypeScript types
  ├── lib/            # Shared libraries
  └── public/         # Static assets
```

## Testing

Run the test suite:

```bash
npm test
```

For test coverage:

```bash
npm run test:coverage
```

## Build

To create a production build:

```bash
npm run build
```

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
