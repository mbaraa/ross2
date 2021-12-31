import MicrosoftLogin from "react-microsoft-login";
import { default as MSLogin } from "../../utils/requests/MicrosoftLogin";
import config from "../../config";
import * as React from "react";
import Button from "@mui/material/Button";

const Login = (): React.ReactElement => {
  // let msal: any = null;
  // try {
  //     msal = JSON.parse(localStorage.getItem("msal") as string);
  // } catch { }

  const [browser, setBrowser] = React.useState<boolean>(false);
  React.useEffect(() => {
    setBrowser(true);
  }, [process.browser]);

  const [msalInstance, onMsalInstanceChange] = React.useState<any | undefined>(
    undefined
  );
  const loginHandler = async (error: any, authData: any, msalInstance: any) => {
    if (!error) {
      console.error(error);
      if (authData) {
        onMsalInstanceChange(msalInstance);
        // localStorage.setItem("msal", JSON.stringify(msalInstance));
      }
    }
    // console.log(authData.account);
    await MSLogin.login(authData);
  };

  return (
    <>
      {browser && (
        <MicrosoftLogin
          clientId={config.msalConfig.auth.clientId}
          authCallback={loginHandler}
          graphScopes={config.msalConfig.scopes}
          redirectUri={config.msalConfig.auth.redirectUri}
          buttonTheme="light_short"
          tenantUrl={config.msalConfig.auth.authority}
        />
      )}
    </>
  );
};

export default Login;
