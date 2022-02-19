import { TextField } from "@mui/material";
import * as React from "react";

export class Point2 {
  x?: number;
  y?: number;
}
export class FieldProps {
  startPosition: Point2;
  width?: number;
  fontSize?: number;
  constructor() {
    this.startPosition = new Point2();
  }
}

interface LabelProps {
  text: string;
}

export function FieldLabel({ text }: LabelProps): React.ReactElement {
  return <label className="font-Ropa text-[16px] text-indigo">{text}</label>;
}

interface Props {
  fieldName: string;
  fieldProps: FieldProps;
  setFieldProps: React.Dispatch<React.SetStateAction<FieldProps>>;
}

const FieldPropsElement = ({
  fieldName,
  fieldProps,
  setFieldProps,
}: Props): React.ReactElement => {
  return (
    <div className="pb-[15px]">
      <label className="font-Ropa text-[20px] text-indigo">{fieldName}</label>
      <div className="grid lg:grid-cols-4 grid-cols-1 pt-[5px]">
        <div className="mr-[10px] mb-[10px]">
          <TextField
            className="w-[100%]"
            variant="outlined"
            value={fieldProps.startPosition.x}
            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
              fieldProps.startPosition.x = Number(event.target.value);
              setFieldProps({ ...fieldProps });
            }}
            label={<FieldLabel text="Position X (px)" />}
            type="number"
          />
        </div>
        <div className="mr-[10px] mb-[10px]">
          <TextField
            className="w-[100%]"
            variant="outlined"
            value={fieldProps.startPosition.y}
            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
              fieldProps.startPosition.y = Number(event.target.value);
              setFieldProps({ ...fieldProps });
            }}
            label={<FieldLabel text="Position Y (px)" />}
            type="number"
          />
        </div>
        <div className="mr-[10px] mb-[10px]">
          <TextField
            className="w-[100%]"
            variant="outlined"
            value={fieldProps.fontSize}
            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
              fieldProps.fontSize = Number(event.target.value);
              setFieldProps({ ...fieldProps });
            }}
            label={<FieldLabel text="Font Size (px)" />}
            type="number"
          />
        </div>
        <div className="mr-[10px]">
          <TextField
            className="w-[100%]"
            variant="outlined"
            value={fieldProps.width}
            onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
              fieldProps.width = Number(event.target.value);
              setFieldProps({ ...fieldProps });
            }}
            label={<FieldLabel text="Field Width (px)" />}
            type="number"
          />
        </div>
      </div>
    </div>
  );
};

export default FieldPropsElement;
