interface Props {
  content: string;
  className: string;
}

const Title = ({ content, className }: Props) => {
  return (
    <div className={`text-ross2 text-[26px] font-[700] + ${className}`}>
      {content}
    </div>
  );
};

export default Title;
