const config = {
    backendAddress: "https://api.ross2.club",
    googleClientID:
      "202655727003-gu3umksjmog90n6oonvfeh79msbe1j1e.apps.googleusercontent.com",
    msalConfig: {
      auth: {
        clientId: "6ab263f3-e5c8-4e89-bbda-ca56faf35977",
        authority:
          "https://login.microsoftonline.com/a6bdeb1e-7724-4165-b796-640034f507ba/",
        redirectUri: "https://ross2.club",
      },
      cache: {
        cacheLocation: "localStorage",
        storeAuthStateInCookie: false,
      },
      scopes: ["openid", "profile", "User.Read"],
    },
  };
  
  export default config;