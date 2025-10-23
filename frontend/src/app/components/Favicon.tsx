'use client';

import Head from 'next/head';

export function Favicon() {
  return (
    <Head>
      <link rel="icon" href="/assets/favicon.ico" />
      <link rel="shortcut icon" href="/assets/favicon.ico" />
      <link rel="apple-touch-icon" href="/assets/favicon.ico" />
    </Head>
  );
} 