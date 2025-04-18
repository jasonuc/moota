type Title = {
    text: string
}

export default function Title({ text }: Title) {
    return (
        <h1 className="text-2xl font-semibold">
            {text}
        </h1>
    )
}