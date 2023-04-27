<svg xmlns="http://www.w3.org/2000/svg" width="114" height="20" role="img" aria-label="coverage: {{printf "%.2f" .Percentage }}%"><title>
        coverage: {{printf "%.2f" .Percentage }}%</title>
    <linearGradient id="s" x2="0" y2="100%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
        <stop offset="1" stop-opacity=".1"/>
    </linearGradient>
    <clipPath id="r">
        <rect width="114" height="20" rx="3" fill="#fff"/>
    </clipPath>
    <g clip-path="url(#r)">
        <rect width="61" height="20" fill="#555"/>
        <rect x="61" width="53" height="20" fill="{{ .Color }}"/>
        <rect width="114" height="20" fill="url(#s)"/>
    </g>
    <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif"
       text-rendering="geometricPrecision" font-size="110">
        <text aria-hidden="true" x="315" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)"
              textLength="510">coverage
        </text>
        <text x="315" y="140" transform="scale(.1)" fill="#fff" textLength="510">coverage</text>
        <text aria-hidden="true" x="865" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)"
              textLength="430">{{printf "%.2f" .Percentage }}%
        </text>
        <text x="865" y="140" transform="scale(.1)" fill="#fff" textLength="430">{{printf "%.2f" .Percentage }}%</text>
    </g>
</svg>
