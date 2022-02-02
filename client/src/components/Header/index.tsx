import * as React from "react";
import { AppBar, IconButton, Toolbar } from "@mui/material";
import AccountCircle from "@mui/icons-material/AccountCircle";
import { Info } from "@mui/icons-material";
import Notifications from "../Notifications";
import { Link } from "react-router-dom";

const Header = (): React.ReactElement => {
  return (
    <>
      <AppBar
        position="fixed"
        elevation={0}
        className="border-b-[1px] border-lwhite font-Ropa"
      >
        <Toolbar className="relative bg-white text-ross2 font-bold text-[1.5em]">
          <Link to="/">
            <img
              src="/logo192.png"
              alt="Ross 2"
              className="h-12 w-12 bg-lwhite-2 border-lwhite-1 border-[0.5px] rounded-full mr-2 cursor-pointer hover:opacity-60"
            />
          </Link>
          <div className="absolute right-[10px]">
            <Notifications />

            <Link to="/profile">
              <IconButton
                size="large"
                aria-label="account of current user"
                aria-controls="menu-appbar"
                aria-haspopup="true"
                color="inherit"
              >
                <AccountCircle />
              </IconButton>
            </Link>

            <Link to="/about">
              <IconButton>
                <Info className="text-ross2" />
              </IconButton>
            </Link>
          </div>
        </Toolbar>
      </AppBar>
    </>
  );
};

export default Header;
