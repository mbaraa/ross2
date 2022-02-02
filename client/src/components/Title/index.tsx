interface Props {
  content: any;
  className: string;
}

const Title = ({ content, className }: Props) => {
  return (
    <div className={`text-ross2 font-Ropa text-[2rem] font-[700] + ${className}`}>
      {content}
    </div>
  );
};

export default Title;
