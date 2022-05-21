import * as React from "react";
import BaseError from "./BaseError";

const NotFound = (): React.ReactElement => {
  return <BaseError errMsg="Page Not Found!" />;
};

export default NotFound;
