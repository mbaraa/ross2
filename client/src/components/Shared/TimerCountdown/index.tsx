import * as React from "react";

interface Props {
  endTimestamp: number;
}

const TimerCountdown = ({ endTimestamp }: Props): React.ReactElement => {
  const [time, setTime] = React.useState({
    days: "",
    hours: "",
    minutes: "",
    seconds: "",
  });

  const _seconds = () => 1000;

  const _minutes = () => {
    return _seconds() * 60;
  };
  const _hours = () => {
    return _minutes() * 60;
  };
  const _days = () => {
    return _hours() * 24;
  };

  const formatNumber = (num: number): string => {
    return num < 10 ? "0" + num.toString(10) : num.toString(10);
  };

  const calcRemainingTime = () => {
    const timer = setInterval(() => {
      const now = new Date();
      const end = new Date(endTimestamp);

      const totalTime = end.getTime() - now.getTime();
      if (totalTime < 0) {
        time.days = "0";
        time.hours = "0";
        time.minutes = "0";
        time.seconds = "0";
        clearInterval(timer);
        setTime({ ...time });
        return;
      }

      const days = Math.floor(totalTime / _days());
      const hours = Math.floor((totalTime % _days()) / _hours());
      const minutes = Math.floor((totalTime % _hours()) / _minutes());
      const seconds = Math.floor((totalTime % _minutes()) / _seconds());
      time.days = formatNumber(days);
      time.hours = formatNumber(hours);
      time.minutes = formatNumber(minutes);
      time.seconds = formatNumber(seconds);
      setTime({ ...time });
    }, 1000);
  };

  calcRemainingTime();

  const getRemainingTime = (): React.ReactNode => {
    return time.days !== "0" &&
      time.hours !== "0" &&
      time.minutes !== "0" &&
      time.seconds !== "0" ? (
      <div>
        <label>{time.days}</label>
        <label className="text-[11px]"> day{time.days !== "01" && "s"}</label>
        {" : "}
        <label>{time.hours}</label>
        <label className="text-[11px]"> hour{time.hours !== "01" && "s"}</label>
        {" : "}
        <label>{time.minutes}</label>
        <label className="text-[11px]">
          {" "}
          minute{time.minutes !== "01" && "s"}
        </label>
      </div>
    ) : (
      "OVER!"
    );
  };

  return <div className="font-[25px]"> {getRemainingTime()} </div>;
};

export default TimerCountdown;
