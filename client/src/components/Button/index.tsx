interface Props {
  content: string;
  onClick: any;
  className: string;
}

/**
  must specify text, border, and hover bg colors in `className`
*/
const Button = ({ content, onClick, className }: Props) => {
  return (
    <div
      onClick={onClick}
      className={`${className} cursor-pointer border-[1px] w-auto inline-block px-[14px] py-[12px] rounded-[8px] text-[13px] font-Ropa font-[500] hover:text-[#fff]`}
    >
      {content}
    </div>
  );
};

export default Button;
