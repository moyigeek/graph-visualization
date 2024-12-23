let myChart = null;

async function fetchData(minCount, viewName) {
    const response = await fetch(`http://localhost:5000/nodes?min_count=${minCount}&view=${viewName}`);
    return await response.json();
}

function processData(data) {
    const nodes = new Map();
    const links = [];
    const degrees = new Map();

    // Calculate node degrees
    data.forEach(item => {
        if (!nodes.has(item.from_package)) {
            nodes.set(item.from_package, { value: 0 });
        }
        if (!nodes.has(item.to_package)) {
            nodes.set(item.to_package, { value: 0 });
        }
        nodes.get(item.to_package).value = item.to_depends;
        nodes.get(item.from_package).value = item.from_depends;

        // Calculate degrees
        degrees.set(item.from_package, (degrees.get(item.from_package) || 0) + 1);
        degrees.set(item.to_package, (degrees.get(item.to_package) || 0) + 1);
    });

    // Create nodes array
    const nodesArray = Array.from(nodes.entries()).map(([name, data]) => ({
        name: name,
        value: data.value,
        degree: degrees.get(name),
        symbolSize: Math.sqrt(data.value+1)*1.5 + 10
    }));

    // Create links array
    data.forEach(item => {
        links.push({
            source: item.from_package,
            target: item.to_package,
            value: item.depends_count
        });
    });

    return { nodes: nodesArray, links: links };
}

function initChart() {
    myChart = echarts.init(document.getElementById('graph'));
    updateGraph(1000, 1); // Default to 'draw_arch'
}

async function updateGraph(minCount, viewNumber) {
    // Map view selection to view names
    let viewName=viewNumber;

    const data = await fetchData(minCount, viewName);
    console.log('Raw data:', data); // Debug point 1

    const graphData = processData(data);
    console.log('Processed graph data:', graphData); // Debug point 2

    // Verify nodes exist
    if (graphData.nodes.length === 0) {
        console.error('No nodes generated from data');
        return;
    }

    const maxNodeSize = 400; // Define the maximum node size

    const option = {
        title: {
            text: 'Package Dependencies Graph'
        },
        tooltip: {
            trigger: 'item',
            formatter: function(params) {
                if (params.dataType === 'node') {
                    return `${params.name}`;
                }
                return `${params.data.source}  ${params.data.target}`;
            }
        },
        series: [{
            type: 'graph',
            layout: 'force',
            data: graphData.nodes,
            links: graphData.links,
            roam: true,
            draggable: true,
            label: {
                show: true,
                fontSize: 12,
                position: 'right'
            },
            force: {
                repulsion: 20000,
                edgeLength: 2000,              
                gravity: 0.1,
                friction: 0.6
            },
            symbolSize:function(value) {
                return Math.min(value + 10, maxNodeSize); // Limit node size
            },
            itemStyle: {
                color: '#1f77b4', // Add default color
                borderWidth: 0,
                borderColor: '#fff'
            },
            edgeSymbol: ['none', 'arrow'],
            edgeSymbolSize: [0, 10],
        }]
    };

    console.log('Chart option:', option); // Debug point 3
    myChart.setOption(option);
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    initChart();
    
    const slider = document.getElementById('minCount');
    const sliderValue = document.getElementById('sliderValue');
    const viewSelect = document.getElementById('viewSelect');
    
    slider.addEventListener('input', (e) => {
        sliderValue.textContent = e.target.value;
    });
    
    slider.addEventListener('change', (e) => {
        updateGraph(parseInt(e.target.value), parseInt(viewSelect.value));
    });

    viewSelect.addEventListener('change', (e) => {
        updateGraph(parseInt(slider.value), parseInt(e.target.value));
    });
});

window.addEventListener('resize', () => {
    myChart?.resize();
});