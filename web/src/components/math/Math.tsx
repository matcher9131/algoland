import katex from "katex";
import "katex/dist/katex.min.css";

type KatexProps = {
    tex: string;
    block?: boolean;
};

export function Math({ tex, block = false } : KatexProps) {
    const html = katex.renderToString(tex, {
        displayMode: block,
        throwOnError: false,
        output: "htmlAndMathml"
    });
    
    return (
        <span
            className={block ? "block my-4" : undefined}
            dangerouslySetInnerHTML={{ __html: html }}
        />
    )
}
