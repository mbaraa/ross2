interface Props {
  content: string;
  onClick: any;
  color: string;
  className?: string;
}

/**
  must specify text, border, and hover bg colors in `className`
*/
const Button = ({ content, onClick, color, className }: Props) => {
  const color0 = (color !== "")? color: "#425CBA";

  return (
    <div
      onClick={onClick}
      className={`${(className !== undefined) ? className: ""} cursor-pointer border-[1px] w-auto inline-block px-[14px] py-[12px] rounded-[8px] font-Ropa font-[500] hover:text-[#fff] border-[${color0}] text-[${color0}] hover:bg-[${color0}] text-center text-[16px]`}
    >
      {content}
    </div>
  );
};

export default Button;
