// import * as React from 'react';
// import {useState, useEffect, useMemo, useCallback} from 'react';
// import {createRoot} from 'react-dom/client';
// import {Map, Source, Layer, LayerProps} from 'react-map-gl/maplibre';
// import ControlPanel from './control-panel';

// export const dataLayer: LayerProps = {
//   id: 'data',
//   type: 'fill',
//   paint: {
//     'fill-color': {
//       type: 'interval',
//       property: 'percentile',
//       stops: [
//         [0, '#3288bd'],
//         [1, '#66c2a5'],
//         [2, '#abdda4'],
//         [3, '#e6f598'],
//         [4, '#ffffbf'],
//         [5, '#fee08b'],
//         [6, '#fdae61'],
//         [7, '#f46d43'],
//         [8, '#d53e4f']
//       ]
//     },
//     'fill-opacity': 0.8
//   }
// };
// export default function MapPage() {
//   const [year, setYear] = useState(2015);
//   const [allData, setAllData] = useState(null);
//   const [hoverInfo, setHoverInfo] = useState(null);

//   useEffect(() => {
//     /* global fetch */
//     fetch(
//       'https://raw.githubusercontent.com/uber/react-map-gl/master/examples/.data/us-income.geojson'
//     )
//       .then(resp => resp.json())
//       .then(json => setAllData(json))
//       .catch(err => console.error('Could not load data', err)); // eslint-disable-line
//   }, []);

//   const onHover = useCallback(event => {
//     const {
//       features,
//       point: {x, y}
//     } = event;
//     const hoveredFeature = features && features[0];

//     // prettier-ignore
//     setHoverInfo(hoveredFeature && {feature: hoveredFeature, x, y});
//   }, []);

//   const data = useMemo(() => {
//     return allData && updatePercentiles(allData, f => f.properties.income[year]);
//   }, [allData, year]);

//   return (
//     <>
//       <Map
//         initialViewState={{
//           latitude: 40,
//           longitude: -100,
//           zoom: 3
//         }}
//         mapStyle="https://basemaps.cartocdn.com/gl/positron-gl-style/style.json"
//         interactiveLayerIds={['data']}
//         onMouseMove={onHover}
//       >
//         <Source type="geojson" data={data}>
//           <Layer {...dataLayer} />
//         </Source>
//       </Map>

//       <ControlPanel year={year} onChange={value => setYear(value)} />
//     </>
//   );
// }

// export function renderToDom(container) {
//   createRoot(container).render(<App />);
// }