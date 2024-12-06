import "./topbar.css";

interface TopBarProps {
  navlinks?: React.ReactNode;
  rightContainer?: React.ReactNode;
}

export const TopBar = ({ navlinks, rightContainer }: TopBarProps) => {
  return (
    <div className="topbar h-12 p-2">
      <div className="container mx-auto">
        <div className="flex flex-row">
          <div>{navlinks}</div>
          <div className="ml-auto">{rightContainer}</div>
        </div>
      </div>
    </div>
  );
};
