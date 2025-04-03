'use client';

import Head from 'next/head';

export function Favicon() {
  return (
    <Head>
      <link rel="icon" href="/assets/UntetherLogo.png" type="image/png" />
      <link rel="shortcut icon" href="/assets/UntetherLogo.png" type="image/png" />
      <link rel="apple-touch-icon" href="/assets/UntetherLogo.png" />
    </Head>
  );
} 