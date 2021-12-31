import "../styles/globals.css";
import type { AppProps } from "next/app";
import Header from "../src/components/Header";
import * as React from "react";
import dynamic from "next/dynamic";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Header></Header>
      <div className="font-[Poppins] absolute left-[0.2rem] top-[4.2em] w-full p-[52px]">
        <Component {...pageProps} />
      </div>
    </>
  );
}

export default MyApp;
