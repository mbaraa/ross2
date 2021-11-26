const config = {
    backendAddress: "http://localhost:8080",
    googleClientID: "",
    msalConfig: {
        auth: {
            clientId: "",
            authority: "https://login.microsoftonline.com/CLIENT_ID",
            redirectUri: "http://localhost:8081",
        },
        cache: {
            cacheLocation: "localStorage",
            storeAuthStateInCookie: false,
        }
    }
};

export default config;
