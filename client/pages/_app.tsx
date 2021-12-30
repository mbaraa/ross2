import "../styles/globals.css";
import type { AppProps } from "next/app";
import Header from "../src/components/Header";
import * as React from "react";
import dynamic from 'next/dynamic'


function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Header></Header>
      <Component {...pageProps} />
    </>
  );
}

export default MyApp;
