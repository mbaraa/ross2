import * as React from "react";
import { Button } from "@mui/material";
import { MdImage } from "react-icons/md";
import { readFile } from "../../../utils";

interface Props {
  /**
   * @maxSize represents image's file size in kilo bytes
   */
  maxSize: number;
  imageFile: File;
  setImageFile: React.Dispatch<React.SetStateAction<File>>;
  className?: string;
}

const ImageUploader = ({
  maxSize,
  imageFile,
  setImageFile,
  className,
}: Props): React.ReactElement => {
  const img = document.getElementById("image-to-upload") as HTMLImageElement;

  React.useEffect(() => {
    if (imageFile.name === "") {
      return;
    }
    (async () => {
      if (imageFile !== null) {
        if (!imageFile?.type.includes("image")) {
          setErrMsg("Whoa... this file is not an image!");
          img.src = "";
          return;
        }

        if (imageFile?.size > maxSize*1000) {
          setErrMsg(`Max file size allowed is ${maxSize}KB`);
          img.src = "";
          return;
        }
        setErrMsg("");
        img.src = `${await readFile(imageFile)}`;
      }
    })();
  }, [imageFile]);

  const [errMsg, setErrMsg] = React.useState("");

  return (
    <div className={`items-center text-center ${className}`}>
      <img
        className={`rounded-[10%] w-[350px] h-[350px] p-[15px] bg-gray-100 relative translate-x-[-50%] left-[50%] ${className}`}
        id="image-to-upload"
        src="/image_upload.png"
        alt=""
      />
      <br />
      <input
        accept="image/*"
        style={{ display: "none" }}
        id="raised-button-file"
        type="file"
        size={5120}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          setImageFile((e.target as any).files[0] as File);
        }}
      />

      <label htmlFor="raised-button-file">
        <Button variant="outlined" component="span" startIcon={<MdImage />}>
          Upload
        </Button>
      </label>

      {errMsg.length > 0 && (
        <label className="text-red-500 font-Ropa text-[20px]">
          <br />
          {errMsg}
        </label>
      )}
    </div>
  );
};

export default ImageUploader;
