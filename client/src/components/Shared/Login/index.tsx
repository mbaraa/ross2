import MicrosoftLogin from "react-microsoft-login";
import { default as MSLogin } from "../../../utils/requests/MicrosoftLogin";
import config from "../../../config";
import * as React from "react";
import { useHistory } from "react-router-dom";

const Login = (): React.ReactElement => {
  const router = useHistory();
  return (
    <div className="font-Ropa">
      <MicrosoftLogin
        clientId={config.msalConfig.auth.clientId}
        authCallback={(error: any, authData: any, msalInstance: any) => {
          (async () => {
            await MSLogin.login(authData);
            router.push("/");
            router.go(0);
          })();
        }}
        graphScopes={config.msalConfig.scopes}
        redirectUri={config.msalConfig.auth.redirectUri}
        buttonTheme="light_short"
        tenantUrl={config.msalConfig.auth.authority}
        prompt="select_account"
      />
    </div>
  );
};

export default Login;
