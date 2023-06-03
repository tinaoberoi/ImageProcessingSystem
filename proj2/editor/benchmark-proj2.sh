#!/bin/bash
#
#SBATCH --mail-user=toberoi@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=project2-time
#SBATCH --output=/home/toberoi/project-2-tinaoberoi/proj2/benchmark/data.txt
#SBATCH --error=./%j.%N.stderr
#SBATCH --chdir=/home/toberoi/project-2-tinaoberoi/proj2/editor
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=03:45:00


module load golang/1.16.2

for i in {1..5}
do
    go run editor.go small
done

for i in {1..5}
do
    go run editor.go mixture
done

for i in {1..5}
do
    go run editor.go big
done

for i in {1..5}
do
    go run editor.go small pipeline 2
done

for i in {1..5}
do
    go run editor.go mixture pipeline 2
done

for i in {1..5}
do
    go run editor.go big pipeline 2
done

for i in {1..5}
do
    go run editor.go small pipeline 4
done

for i in {1..5}
do
    go run editor.go mixture pipeline 4
done

for i in {1..5}
do
    go run editor.go big pipeline 4
done

for i in {1..5}
do
    go run editor.go small pipeline 6
done

for i in {1..5}
do
    go run editor.go mixture pipeline 6
done

for i in {1..5}
do
    go run editor.go big pipeline 6
done

for i in {1..5}
do
    go run editor.go small pipeline 8
done

for i in {1..5}
do
    go run editor.go mixture pipeline 8
done

for i in {1..5}
do
    go run editor.go big pipeline 8
done

for i in {1..5}
do
    go run editor.go small pipeline 12
done

for i in {1..5}
do
    go run editor.go mixture pipeline 12
done

for i in {1..5}
do
    go run editor.go big pipeline 12
done

for i in {1..5}
do
    go run editor.go small bsp 2
done

for i in {1..5}
do
    go run editor.go mixture bsp 2
done

for i in {1..5}
do
    go run editor.go big bsp 2
done

for i in {1..5}
do
    go run editor.go small bsp 4
done

for i in {1..5}
do
    go run editor.go mixture bsp 4
done

for i in {1..5}
do
    go run editor.go big bsp 4
done

for i in {1..5}
do
    go run editor.go small bsp 6
done

for i in {1..5}
do
    go run editor.go mixture bsp 6
done

for i in {1..5}
do
    go run editor.go big bsp 6
done

for i in {1..5}
do
    go run editor.go small bsp 8
done

for i in {1..5}
do
    go run editor.go mixture bsp 8
done

for i in {1..5}
do
    go run editor.go big bsp 8
done

for i in {1..5}
do
    go run editor.go small bsp 12
done

for i in {1..5}
do
    go run editor.go mixture bsp 12
done

for i in {1..5}
do
    go run editor.go big bsp 12
done

python3 plot.py