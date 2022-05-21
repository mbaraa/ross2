import * as React from "react";
import Title from "../Title";

interface Props {
  errMsg?: string;
}

const BaseError = ({ errMsg }: Props): React.ReactElement => {
  return (
    <Title
      className="text-red-600 absolute top-[50%] left-[50%] translate-x-[-50%] translate-y-[-50%] bg-white w-[100vw] h-[100%] z-[50] ml-[25px]"
      content={errMsg ?? "Something went wrong..."}
    />
  );
};

export default BaseError;
