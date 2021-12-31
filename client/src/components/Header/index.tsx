import * as React from "react";
import {
  AppBar,
  Badge,
  IconButton,
  Menu,
  MenuItem,
  Toolbar,
} from "@mui/material";
import AccountCircle from "@mui/icons-material/AccountCircle";
import NotificationsIcon from "@mui/icons-material/Notifications";
import { Info } from "@mui/icons-material";
import Notification from "../../models/Notification";
import NotificationRequests from "../../utils/requests/NotificationRequests";
import Team from "../../models/Team";
import ContestantRequests from "../../utils/requests/ContestantRequests";
import MicrosoftLogin from "../../utils/requests/MicrosoftLogin";
import User from "../../models/User";
import Login from "../Login";
import { useRouter } from "next/router";

interface Props {
  user: User;
}

const ProfileMenu = ({ user }: Props): React.ReactElement => {
  const router = useRouter();

  const profileMenu = [
    {
      id: 1,
      name: "Profile",
      action: () => {
        router.push("/profile");
      },
    },
    {
      id: 2,
      name: "Team",
      action: () => {
        router.push("/team");
      },
    },
    // { name: "Logout", action: () => {  } }
  ];

  return !user ? (
    <Login />
  ) : (
    <>
      {profileMenu.map((m: any) => (
        <MenuItem key={m.id} onClick={() => m.action()}>
          {m.name}
        </MenuItem>
      ))}
    </>
  );
};

const Header = (): React.ReactElement => {
  const profileMenu = ["Profile", "Team", "Logout"];
  const [menu, setMenu] = React.useState<null | HTMLElement>(null);
  const toggleProfileMenu = (event: React.MouseEvent<HTMLElement>) => {
    setMenu(event.currentTarget);
  };
  const closeMenu = () => {
    setMenu(null);
  };

  const [nots, setNots] = React.useState<Notification[]>([]);
  React.useEffect(() => {
    setNotifications();
  }, []);
  const setNotifications = async () => {
    const n = await NotificationRequests.getNotifications();
    setNots(n);
  };

  const [team, setTeam] = React.useState<Team>(new Team());
  React.useEffect(() => {
    setteam();
  }, []);

  const setteam = async () => {
    const t = await ContestantRequests.getTeam();
    setTeam(t);
  };

  const [user, setUser] = React.useState<User>(new User());
  React.useEffect(() => {
    login();
  }, []);
  const login = async () => {
    const u = await MicrosoftLogin.loginWithToken();
    setUser(u);
  };

  return (
    <>
      <AppBar
        position="fixed"
        elevation={0}
        className="border-b-[1px] border-lwhite"
      >
        <Toolbar className="relative bg-white text-ross2 font-bold text-[1.5em]">
          <a href="/">
            <img
              src="/logo192.png"
              alt="Ross 2"
              className="h-12 w-12 bg-lwhite-2 border-lwhite-1 border-[0.5px] rounded-full mr-2 cursor-pointer hover:opacity-60"
            />
          </a>
          <div className="absolute right-[10px]">
            <IconButton size="large" aria-label="notifications">
              <Badge badgeContent={nots.length} color="error">
                <NotificationsIcon className="text-ross2" />
              </Badge>
            </IconButton>

            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={toggleProfileMenu}
              color="inherit"
            >
              <AccountCircle />
            </IconButton>
            <Menu
              anchorEl={menu}
              anchorOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              keepMounted
              transformOrigin={{
                vertical: "top",
                horizontal: "right",
              }}
              open={Boolean(menu)}
              onClose={closeMenu}
            >
              <ProfileMenu user={user} />
            </Menu>

            <a href="/about">
              <IconButton>
                <Info className="text-ross2" />
              </IconButton>
            </a>
          </div>
        </Toolbar>
      </AppBar>
    </>
  );
};

export default Header;
