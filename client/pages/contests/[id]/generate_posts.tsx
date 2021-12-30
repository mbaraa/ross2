import { useRouter } from "next/router";

const GeneratePosts = () => {
  const router = useRouter();

  return <div>GeneratePosts {router.query["id"]}</div>;
};

export default GeneratePosts;
