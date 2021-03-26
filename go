set -e

../tools/acme -o new/miner2021 -r miner2021.txt --vicelabels miner2021.sym miner2021.a
../tools/acme -o new/loader -r loader.txt --vicelabels loader.sym loader.a
sort -o miner2021.tmp miner2021.sym
sort -o loader.tmp loader.sym
uniq miner2021.tmp miner2021.sym
uniq loader.tmp loader.sym
rm miner2021.tmp
rm loader.tmp

echo
# Create new SSD file with the appropriate files
../tools/bbcim -new -type ADFS NEWMINER.ssd
../tools/bbcim -boot NEWMINER.ssd EXEC
../tools/bbcim -a NEWMINER.ssd new/!BOOT new/MANIC new/LOADER new/MINER2021

echo
if [ $USER == "tobynelson" ];
then
    # Open SSD in b2
    DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

    open -a 'b2 Debug'
    sleep 1
    # curl -G 'http://localhost:48075/reset/b2' --data-urlencode "config=Master 128 (MOS 3.50)"
    curl -H 'Content-Type:application/binary' --upload-file "$DIR/NEWMINER.ssd" 'http://localhost:48075/run/b2?name=NEWMINER.ssd'

else
    # Open SSD in BeebEm
    open NEWMINER.ssd
fi
