import pandas as pd
import numpy as np
import yaml as yml
import sys

####
# The script will print all inconsistencies between the results and the input.
# You will need to have sufficient data in the database for good results.
# To call the script:
#
# python3 check_perm_error.py <path_to_xlsx_with_db_results> <path_to_yaml>
#
####

# Load xlsx into memory
def load_xlsx(file_name):
    xlsx_df = pd.read_excel(file_name, sheet_name='Sheet1')

    # Remove the first 2 columns
    xlsx_df = xlsx_df.drop(xlsx_df.columns[0], axis=1)
    xlsx_df = xlsx_df.drop(xlsx_df.columns[0], axis=1)

    return xlsx_df

def format_columns(xlsx):
  # Create a dictionary between the column names and unique values
  col_dict = {}

  # Iterate through each column
  for col in xlsx.columns:
    # Add each row to a list
    row_list = []
    for row in xlsx[col]:
      row_list.append(row)

    # Sort the list and remove duplicates and stringify values
    row_list = list(set([str(val) for val in row_list]))
    row_list.sort()

    # Add the column name and the list to the dictionary
    col_dict[col] = row_list

  return col_dict

# Load yaml into memory
def load_yaml(file_name):
  with open(file_name, 'r') as file:
    job_yaml = yml.safe_load(file)

  yaml_dict = {}

  for param in job_yaml['params']:
    if 'only' in param:
      yaml_dict[param['name']] = list(set([str(val) for val in param['only']]))
      yaml_dict[param['name']].sort()
    else:
      yaml_dict[param['name']] = str(param['default'])

  return yaml_dict

def main():
  # Load xlsx
  xlsx = load_xlsx(sys.argv[1])

  # Load yaml
  yaml_dict = load_yaml(sys.argv[2])

  # Sort and unique columns
  xlsx_dict = format_columns(xlsx)

  # Iterate over the yaml_dict and compare it to xlsx_dict
  for param_name in yaml_dict:
    if param_name in yaml_dict and param_name in xlsx_dict:
      if yaml_dict[param_name] != xlsx_dict[param_name]:
        print(param_name, ": yaml =", yaml_dict[param_name], " xlsx =", xlsx_dict[param_name])
    else:
      if not param_name in yaml_dict:
        print(param_name, "missing in yaml")
      if not param_name in xlsx_dict:
        print(param_name, "missing in xlsx")

if __name__ == '__main__':
  main()