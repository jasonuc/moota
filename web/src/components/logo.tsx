export function Logo() {
  return (
    <img
      src="/moota.png"
      alt="Moota"
      width={45}
      height={45}
      draggable={false}
    />
  );
}

export function LogoWithText() {
  return (
    <div className="flex items-center justify-center space-x-2 grow-0 w-fit">
      <Logo />
      <h1 className="font-bold text-3xl">Moota</h1>
    </div>
  );
}
