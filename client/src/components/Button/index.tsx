interface Props {
  content: string;
  onClick: any;
  className: string;
}

const Button = ({ content, onClick, className }: Props) => {
  return (
    <div
      onClick={onClick}
      className={`cursor-pointer border-[1px] border-[#425CBA] w-auto inline-block px-[14px] py-[12px] rounded-[8px] text-[13px] text-[#425CBA] font-Ropa font-[500] hover:bg-[#425CBA] hover:text-[#fff] ${className}`}
    >
      {content}
    </div>
  );
};

export default Button;
