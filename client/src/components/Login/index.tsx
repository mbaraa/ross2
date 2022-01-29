import MicrosoftLogin from "react-microsoft-login";
import { default as MSLogin } from "../../utils/requests/MicrosoftLogin";
import config from "../../config";
import * as React from "react";

const Login = (): React.ReactElement => {
  return (
    <MicrosoftLogin
      clientId={config.msalConfig.auth.clientId}
      authCallback={(error: any, authData: any, msalInstance: any) => {
        (async () => {
          await MSLogin.login(authData);
          window.location.reload();
        })();
      }}
      graphScopes={config.msalConfig.scopes}
      redirectUri={config.msalConfig.auth.redirectUri}
      buttonTheme="light_short"
      tenantUrl={config.msalConfig.auth.authority}
      prompt="select_account"
    />
  );
};

export default Login;
