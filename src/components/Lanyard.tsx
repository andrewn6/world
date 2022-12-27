import { useLanyardWS } from "use-lanyard";

export default function Lanyard({ discordId }: { discordId: string }) {
  const data = useLanyardWS("804052488235647017");

  const activity = data?.activities;

  if (!activity) {
    return <></>;
  }

  return (
    <p className="flex space-x-1">
      <p>•</p>{" "}
      <p>
        Currently doing {activity.song}
       </p>
    </p>
)}
