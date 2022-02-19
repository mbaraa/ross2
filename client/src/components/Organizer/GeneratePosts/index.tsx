import * as React from "react";
import Contest, { ParticipationConditions } from "../../../models/Contest";
import { Button, FormControlLabel, Switch, TextField } from "@mui/material";
import OrganizerRequests from "../../../utils/requests/OrganizerRequests";
import { MdImage } from "react-icons/md";
import { GiGears } from "react-icons/gi";
import ImageUploader from "../../Shared/ImageUploader";
import FieldPropsElement, { FieldProps, FieldLabel } from "./FieldPropsElement";
import { readFile } from "../../../utils";
interface Props {
  contest: Contest;
}

const GeneratePosts = ({ contest }: Props): React.ReactElement => {
  const [useSample, setUseSample] = React.useState(false);
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUseSample(event.target.checked);
  };

  const [templateImage, setTemplateImage] = React.useState<File>(new File([], ""));
  const [teamNumberProps, setTeamNumberProps] = React.useState(
    new FieldProps()
  );
  const [teamNameProps, setTeamNameProps] = React.useState(new FieldProps());

  const membersProps: any = [];

  for (
    let i = 0;
    i <
    (((contest as Contest).participation_conditions as ParticipationConditions)
      .max_team_members as number);
    i++
  ) {
    membersProps.push({
      i: i,
      fp: new FieldProps(),
    });
  }

  const [fp, setFP] = React.useState(membersProps);

  const checkImageFile = (): boolean => {
            if (!useSample && (templateImage === undefined || templateImage?.type.indexOf("png") === -1)) {
                window.alert("select file of image/png type!");
                setTemplateImage(new File([], ""));
                // templateImage = null;
                return false;
            }
            return true;
        }

  const generatePosts = () => {
    if (!checkImageFile()) {
      return;
  }

    const theRealMembersNamesProps = new Array<FieldProps>();

    fp.forEach((fpi: any) => {
      theRealMembersNamesProps.push(fpi.fp);
    });

    (async () => {
      const templateImageB64 = await readFile(templateImage as File) as string;
      const zipFile = await OrganizerRequests.generateTeamsPosts({
        "contest": contest,
        "teamNameProps": teamNameProps,
        "teamOrderProps": teamNumberProps,
        "membersNamesProps": theRealMembersNamesProps,
        "baseImage": useSample? "": templateImageB64?.substring(templateImageB64?.indexOf(",")+1),
      });

      const a = document.createElement("a");
      a.href = `data:application/zip;base64,${zipFile}`;
      a.download = `${contest.name}'s_teams_posts.zip`;
      a.click();
    })();

  };

  return (
    <>
      {/* sample images */}
      <div className="grid grid-cols-1 md:grid-cols-4">
        <div>
          <Button
            variant="outlined"
            color="secondary"
            startIcon={<MdImage size={15} />}
            onClick={() => window.open("/team_post_template.png", "_blank")}
          >
            <label className="normal-case font-Ropa text-[15px] cursor-pointer">
              Download sample template image
            </label>
          </Button>
        </div>
        <div>
          <Button
            variant="outlined"
            color="info"
            startIcon={<MdImage size={15} />}
            onClick={() => window.open("/team_post_sample.png", "_blank")}
          >
            <label className="normal-case font-Ropa text-[15px] cursor-pointer">
              Download sample filled image
            </label>
          </Button>
        </div>
      </div>
      {/*  */}
      <div className="grid grid-cols-1 sm:grid-cols-4 pt-[10px]">
        <div>
          <FormControlLabel
            className="py-[10px]"
            control={
              <Switch
                color="info"
                checked={useSample}
                onChange={handleChange}
                inputProps={{ "aria-label": "controlled" }}
              />
            }
            label={
              <label className="font-Ropa text-indigo font-[18px]">
                Use Sample Image (for 3 members per team)
              </label>
            }
          />
        </div>
        <div>
          <Button
            variant="outlined"
            color="error"
            startIcon={<GiGears size={15} />}
            onClick={generatePosts}
          >
            <label className="normal-case font-Ropa text-[15px] cursor-pointer">
              Generate Posts
            </label>
          </Button>
        </div>
      </div>
      {/* post props */}
      {!useSample && (
        <div className="grid grid-cols-1 lg:grid-cols-2">
          {/* left side */}
          <div>
            <ImageUploader
              maxSize={5120}
              imageFile={templateImage}
              setImageFile={setTemplateImage}
              className="xl:w-[600px] rounded-none"
            />
          </div>
          {/* right side */}
          <div>
            <FieldPropsElement
              fieldName="Team Number Props:"
              fieldProps={teamNumberProps}
              setFieldProps={setTeamNumberProps}
            />
            <FieldPropsElement
              fieldName="Team Name Props:"
              fieldProps={teamNameProps}
              setFieldProps={setTeamNameProps}
            />
            {/* well I couldn't :) */}
            {fp.map((mp: any) => (
              <div className="pb-[15px]" key={mp.i}>
                <label className="font-Ropa text-[20px] text-indigo">{`Member #${
                  mp.i + 1
                } Name Props:`}</label>
                <div className="grid lg:grid-cols-4 grid-cols-1 pt-[5px]">
                  <div className="mr-[10px] mb-[10px]">
                    <TextField
                      className="w-[100%]"
                      variant="outlined"
                      value={fp[mp.i].fp.startPosition.x}
                      onChange={(
                        event: React.ChangeEvent<HTMLInputElement>
                      ) => {
                        fp[mp.i].fp.startPosition.x = Number(
                          event.target.value
                        );
                        setFP(fp);
                      }}
                      label={<FieldLabel text="Position X (px)" />}
                      type="number"
                    />
                  </div>
                  <div className="mr-[10px] mb-[10px]">
                    <TextField
                      className="w-[100%]"
                      variant="outlined"
                      value={fp[mp.i].fp.startPosition.y}
                      onChange={(
                        event: React.ChangeEvent<HTMLInputElement>
                      ) => {
                        fp[mp.i].fp.startPosition.y = Number(
                          event.target.value
                        );
                        setFP(fp);
                      }}
                      label={<FieldLabel text="Position Y (px)" />}
                      type="number"
                    />
                  </div>
                  <div className="mr-[10px] mb-[10px]">
                    <TextField
                      className="w-[100%]"
                      variant="outlined"
                      value={fp[mp.i].fp.fontSize}
                      onChange={(
                        event: React.ChangeEvent<HTMLInputElement>
                      ) => {
                        fp[mp.i].fp.fontSize = Number(event.target.value);
                        setFP(fp);
                      }}
                      label={<FieldLabel text="Font Size (px)" />}
                      type="number"
                    />
                  </div>
                  <div className="mr-[10px]">
                    <TextField
                      className="w-[100%]"
                      variant="outlined"
                      value={fp[mp.i].fp.width}
                      onChange={(
                        event: React.ChangeEvent<HTMLInputElement>
                      ) => {
                        fp[mp.i].fp.width = Number(event.target.value);
                        setFP(fp);
                      }}
                      label={<FieldLabel text="Field Width (px)" />}
                      type="number"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </>
  );
};

export default GeneratePosts;
