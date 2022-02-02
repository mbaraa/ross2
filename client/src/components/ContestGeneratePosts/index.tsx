import { useHistory, useParams } from "react-router-dom";
import * as React from "react";
import Contest from "../../../src/models/Contest";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Input,
} from "@mui/material";
import OrganizerRequests from "../../../src/utils/requests/OrganizerRequests";

class Point2 {
  x?: number;
  y?: number;
}
class FieldProps {
  startPosition: Point2;
  width?: number;
  fontSize?: number;
  constructor() {
    this.startPosition = new Point2();
  }
}

const ContestGeneratePosts = (): React.ReactElement => {
  const router = useHistory();
  const { id }: any = useParams();

  const [contest, setCont] = React.useState<Contest>(new Contest());
  React.useEffect(() => {
    (async () => {
      const c = await Contest.getContestFromServer(
        parseInt(id as string)
      );
      setCont(c);
    })();
  }, []);
  
  const [useSampleImage, setUseSampleImage] = React.useState<boolean>(false);
  const [loading, setLoading] = React.useState<boolean>(false);

  const [teamNameProps, setTeamNameProps] = React.useState<FieldProps>(
    new FieldProps()
  );
  const [teamNumberProps, setTeamNumberProps] = React.useState<FieldProps>(
    new FieldProps()
  );
  const [membersProps, setMembersProps] = React.useState<any[]>([]);
  const [imageB64, setImageB64] = React.useState<string>("");
  const [templateFile, setTemplateFile] = React.useState<any>(undefined);

  const [postsGenDialog, setPostsGenDialog] = React.useState<boolean>(false);

  // const openPostsGenDialog = (contest: Contest) => {
    // setCont(contest);

    const mb = new Array<FieldProps>(
      contest.participation_conditions?.max_team_members as number
    ).fill({width: 0, fontSize: 0, startPosition: new Point2()} as FieldProps);
    // class Foo extends FieldProps {
    //   id: number;
    //   constructor(id: number) {
    //     super();
    //     this.id = id;
    //   }
    // }

    // for (let i = 0; i < (contest.participation_conditions?.max_team_members as number); i++) {
    //   mb.push(new Foo(i));
    // }

    console.log("mb", mb);
    setMembersProps(mb);
    setPostsGenDialog(true);
  // };

  const closePostsGenDialog = () => {
    setPostsGenDialog(false);
  };


  const checkImageFile = (): boolean => {
    if (
      !useSampleImage &&
      (templateFile === undefined || templateFile.type.indexOf("png") == -1)
    ) {
      window.alert("select file of image/png type!");
      setTemplateFile(undefined);
      return false;
    }
    return true;
  };

  const selectFile = (file: any) => {
    setTemplateFile(file.target.files[0]);
    checkImageFile();
  };

  const readFile = async (): Promise<string | ArrayBuffer | null> => {
    let res: string | ArrayBuffer | null = "";
    // 🙉🙊🙈 if it works it ain't stupid
    const toBase64 = () =>
      new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(templateFile);
        reader.onload = () => {
          resolve(reader.result);
          res = reader.result;
          return res;
        };
        reader.onerror = (error) => reject(error);
      });
    await toBase64();
    return res;
  };

  const generatePosts = async (contest: Contest) => {
    if (!checkImageFile()) {
      return;
    }
    setLoading(true);
    const templateImageB64 = useSampleImage ? "" : await readFile();
    const zipFile = await OrganizerRequests.generateTeamsPosts({
      contest: contest,
      teamNameProps: teamNameProps,
      teamOrderProps: teamNumberProps,
      membersNamesProps: membersProps,
      baseImage: useSampleImage
        ? ""
        : (templateImageB64 as string).substring(
            (templateImageB64 as string).indexOf(",") + 1
          ),
    });

    setLoading(false);

    const f = document.createElement("a");
    f.href = `data:application/zip;base64,${zipFile}`;
    f.download = `${contest.name}'s teams posts.zip`;
    f.click();
  };

  const [memberNumber, setMemberNumber] = React.useState<number>(0);

  return (
    <>
      {/* <Dialog open={postsGenDialog} onClose={closePostsGenDialog}> */}
        <p className="text-[Poppins] text-indigo font-[600] font-Ropa">
          Generate Teams Posts
        </p>
        {/* <DialogContent className="text-[Poppins]"> */}
          <p>Team number props:</p>
          <Input
            type="number"
            required
            placeholder="Start position X"
            value={teamNameProps.startPosition.x}
          />
          <br />
          <Input
            type="number"
            required
            placeholder="Start position Y"
            value={teamNameProps.startPosition.y}
          />
          <br />
          <Input
            type="number"
            required
            placeholder="Font Size"
            value={teamNameProps.fontSize}
          />
          <br />
          <Input
            type="number"
            required
            placeholder="Width"
            value={teamNameProps.width}
          />

          {membersProps.map((mp: FieldProps) => {
            setMemberNumber(memberNumber + 1);
            return (
              <div key={Math.random()}>
                <p>Member #{memberNumber} name props:</p>
                <Input
                  type="number"
                  required
                  placeholder="Start position X"
                  value={mp.startPosition.x}
                />
                <br />
                <Input
                  type="number"
                  required
                  placeholder="Start position Y"
                  value={mp.startPosition.y}
                />
                <br />
                <Input
                  type="number"
                  required
                  placeholder="Font Size"
                  value={mp.fontSize}
                />
                <br />
                <Input
                  type="number"
                  required
                  placeholder="Width"
                  value={mp.width}
                />
                <br />
              </div>
            );
          })}
        {/* </DialogContent> */}
        {/* <DialogActions> */}
          <Button onClick={closePostsGenDialog}>Cancel</Button>
          <Button onClick={closePostsGenDialog}>Subscribe</Button>
        {/* </DialogActions> */}
      {/* </Dialog> */}
    </>
  );
};

export default ContestGeneratePosts;