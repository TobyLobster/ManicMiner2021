set -e

# ../tools/tokenizer2 manic.txt >new/manic
../tools/acme -o new/miner1 -r miner1.txt --vicelabels miner1.sym miner1.a
../tools/acme -o new/miner2 -r miner2.txt --vicelabels miner2.sym miner2.a
../tools/acme -o new/miner3 -r miner3.txt --vicelabels miner3.sym miner3.a
../tools/acme -o new/miner4 -r miner4.txt --vicelabels miner4.sym miner4.a
if [ "$1" != "skip" ]; then
    echo Checking diffs...
    diff new/manic  ./original/Disc011-ManicMiner.ssd.d/MANIC
    diff new/miner1 ./original/Disc011-ManicMiner.ssd.d/MINER1.bin
    diff new/miner2 ./original/Disc011-ManicMiner.ssd.d/MINER2.bin
    diff new/miner3 ./original/Disc011-ManicMiner.ssd.d/MINER3.bin
    diff new/miner4 ./original/Disc011-ManicMiner.ssd.d/MINER4.bin
else
    echo Skipping diffs...
fi
sort -o miner1.tmp miner1.sym
sort -o miner2.tmp miner2.sym
sort -o miner3.tmp miner3.sym
sort -o miner4.tmp miner4.sym
uniq miner1.tmp miner1.sym
uniq miner2.tmp miner2.sym
uniq miner3.tmp miner3.sym
uniq miner4.tmp miner4.sym
rm miner1.tmp
rm miner2.tmp
rm miner3.tmp
rm miner4.tmp

echo
# Create new SSD file with the appropriate files
../tools/bbcim -new -type ADFS NEWMINER.ssd
../tools/bbcim -boot NEWMINER.ssd EXEC
../tools/bbcim -a NEWMINER.ssd new/!BOOT new/MANIC new/MINER1 new/MINER2 new/MINER3 new/MINER4

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
