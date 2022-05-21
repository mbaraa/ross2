import * as React from "react";
import BaseError from "./BaseError";

const YouCantDoThat = (): React.ReactElement => {
  return <BaseError errMsg="You Can't Do That!" />;
};

export default YouCantDoThat;
