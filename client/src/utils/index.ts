export function getLocaleTime(time: number): string {
  return new Date(time).toLocaleTimeString("en-US", {
    hour12: true,
    hour: "2-digit",
    minute: "2-digit",
    year: "numeric",
    month: "short",
    day: "2-digit",
    weekday: "short",
  });
}

export function formatDuration(minutes: number): string {
  const hours = (minutes / 60) % 24;
  const minutes1 = minutes % 60;

  return `${hours} hours${minutes1 !== 0 ? ` & ${minutes1} minutes` : ""}`;
}

export async function readFile(
  file: File
): Promise<string | ArrayBuffer | null> {
  let res: string | ArrayBuffer | null = "";
  // 🙉🙊🙈 if it works it ain't stupid
  const toBase64 = () =>
    new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        resolve(reader.result);
        res = reader.result;
        return res;
      };
      reader.onerror = (err) => reject(err);
    });
  await toBase64();

  return res;
}
