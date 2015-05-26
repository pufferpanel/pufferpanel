#/bin/bash
version=0.1
divider="-----"
installPath="PufferPanel"

echo "PufferPanel Language Builder - v${version}"
echo $divider

while getopts "h?p:" opt; do
    case "$opt" in
    p)  
        installPath=$OPTARG
        ;;
    h)
        echo "PufferPanel Language Builder - v${version}"
        echo "Optional parameters: "
        echo "-p        | Path to PufferPanel installation (i.e PufferPanel)"
        exit;
        ;;
    esac
done

shift $((OPTIND-1))

[ "$1" = "--" ] && shift

echo "Cleaning up old language files"
rawPath="${installPath}/app/languages/raw/"
outputPath="${installPath}/app/languages/"
rm -f ${outputPath}*.json

echo $divider
echo "Generating language files"

cd $rawPath
echo -n "    "
for f in *.txt
do
    filename=$(basename $f .txt)
    echo "Building ${filename}"
    php << EOF
    <?php 
        \$content = file('${filename}.txt');
        \$json = array();
        foreach(\$content as \$line => \$string){
            list(\$id, \$lang) = explode(",", \$string, 2);
            \$json = array_merge(\$json, array(strtolower(\$id) => trim(\$lang)));
        }
        \$fp = fopen('../${filename}.json', 'w+');
        fwrite(\$fp, json_encode(\$json, JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT));
        fclose(\$fp);
    ?>
EOF
    
done

echo
echo $divider
echo "Building complete"
cd $(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)