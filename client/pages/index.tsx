import { Router, useRouter } from "next/router";
import { useEffect } from "react";
import Table from "../src/components/Table";

const Contests = () => {
  const router = useRouter();

  useEffect(() => {
    router.push("/contests");
  }, []);

  return <div>Loading...</div>;
};

export default Contests;
