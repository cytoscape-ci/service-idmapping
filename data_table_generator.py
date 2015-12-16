# coding: utf-8

# Source data set builder for in-memory data table used in ID Mapper.
import pandas as pd
import json, urllib.request

# Location of original data sets
NCBI_FTP = "ftp://ftp.ncbi.nih.gov/gene/DATA/GENE_INFO/"
UNI_FTP = "ftp://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/idmapping/by_organism/"

# From NCBI Gene
NCBI_SOURCES = {
    "HUMAN": NCBI_FTP + "Mammalia/Homo_sapiens.gene_info.gz",
    "YEAST": NCBI_FTP + "Fungi/Saccharomyces_cerevisiae.gene_info.gz",
    "MOUSE": NCBI_FTP + "Mammalia/Mus_musculus.gene_info.gz",
    "FLY": NCBI_FTP + "Invertebrates/Drosophila_melanogaster.gene_info.gz"
}

# Uniprot ID Mapping
UNIPROT_SOURCES = {
    "HUMAN": UNI_FTP + "HUMAN_9606_idmapping_selected.tab.gz",
    "MOUSE": UNI_FTP + "MOUSE_10090_idmapping_selected.tab.gz",
    "YEAST": UNI_FTP + "YEAST_559292_idmapping_selected.tab.gz",
    "FLY": UNI_FTP + "DROME_7227_idmapping_selected.tab.gz"
}

# Get all data sets from the FTP server
NCBI_COLUMNS = ["tax_id", "GeneID", "Symbol", "LocusTag", "Synonyms", "dbXrefs", "chromosome", "map_location",
                "description", "type_of_gene", "Symbol_from_nomenclature_authority", "Full_name_from_nomenclature_authority",
                "Nomenclature_status", "Other_designations", "Modification_date"]

UNIPROT_COLUMNS = ['UniProtKB-AC','UniProtKB-ID','GeneID','RefSeq','GI','PDB','GO','UniRef100','UniRef90',
    'UniRef50','UniParc','PIR','NCBI-taxon','MIM','UniGene','PubMed','EMBL','EMBL-CDS',
    'Ensembl','Ensembl_TRS','Ensembl_PRO','Additional PubMed']

## Load NCBI data first.
ncbi_map = {}

for key in NCBI_SOURCES:
    print("Downloading from NCBI FTP server: " + key + "...")    
    local_filename, headers = urllib.request.urlretrieve(NCBI_SOURCES[key])
    ncbi_map[key] = pd.read_csv(local_filename, sep='\t', low_memory=False, 
                             names=NCBI_COLUMNS, comment='#', compression="gzip")

# Load UNIPROT data next...
uniprot_map = {}

for key in UNIPROT_SOURCES:
    print("Downloading from Uniprot FTP server: " + key + "...")    
    local_filename, headers = urllib.request.urlretrieve(UNIPROT_SOURCES[key])
    uniprot_map[key] = pd.read_csv(local_filename, 
                     sep="\t", names=UNIPROT_COLUMNS, low_memory=False, compression="gzip")

# Create mapping files
for key in ncbi_map:
    ncbi_gene_info = ncbi_map[key]
    ncbi_subset = ncbi_gene_info[["GeneID", "Symbol", "Full_name_from_nomenclature_authority"]].astype(str)
    
    # Merge and create new table
    merged = pd.merge(uniprot_map[key], ncbi_subset , left_on="GeneID", right_on="GeneID", how="outer")
    # Drop unnecessary columns
    df_final = merged.drop("Additional PubMed", 1)
    df_final = df_final.drop("NCBI-taxon", 1)
    df_final = df_final.drop("PubMed", 1)
    
    # Create one mapping file / species
    df_final.to_csv("./data/idmapping_" + key.lower() +".tsv", sep='\t', index=False)
